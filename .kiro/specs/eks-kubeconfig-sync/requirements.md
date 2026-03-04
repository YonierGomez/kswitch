# Documento de Requisitos

## Introducción

Nuevo subcomando `ksw eks kubeconfig` que sincroniza automáticamente el kubeconfig local con los clústeres EKS disponibles en las cuentas AWS del usuario. Lee los perfiles configurados en `~/.aws/config`, descubre los clústeres EKS en cada perfil/región y agrega al kubeconfig aquellos que aún no existen, evitando duplicados.

## Glosario

- **KSW**: CLI interactivo para cambio de contextos Kubernetes. Binario único escrito en Go.
- **AWS_Config_File**: Archivo `~/.aws/config` que contiene los perfiles AWS del usuario con sus regiones asociadas.
- **Perfil_AWS**: Sección `[profile <nombre>]` dentro del AWS_Config_File que define credenciales y región para una cuenta AWS.
- **Clúster_EKS**: Clúster de Amazon Elastic Kubernetes Service asociado a un Perfil_AWS y una región.
- **Kubeconfig**: Archivo de configuración de Kubernetes (por defecto `~/.kube/config`) que contiene contextos, clústeres y credenciales.
- **Contexto_Kubeconfig**: Entrada dentro del Kubeconfig que referencia un clúster y credenciales específicas.
- **Sincronizador_EKS**: Componente de KSW responsable de descubrir clústeres EKS y agregarlos al Kubeconfig.
- **AWS_CLI**: Herramienta de línea de comandos de AWS, requerida para ejecutar `aws eks update-kubeconfig`.

## Requisitos

### Requisito 1: Lectura de perfiles AWS

**User Story:** Como operador de Kubernetes, quiero que KSW lea mis perfiles AWS configurados, para que pueda descubrir clústeres EKS en todas mis cuentas sin configuración adicional.

#### Criterios de Aceptación

1. WHEN el usuario ejecuta `ksw eks kubeconfig`, THE Sincronizador_EKS SHALL leer todos los perfiles definidos en el AWS_Config_File (`~/.aws/config`).
2. WHEN un Perfil_AWS tiene una región configurada, THE Sincronizador_EKS SHALL utilizar esa región para la búsqueda de clústeres EKS en ese perfil.
3. WHEN un Perfil_AWS no tiene región configurada, THE Sincronizador_EKS SHALL utilizar la región por defecto del AWS_CLI (`aws configure get region`) o `us-east-1` como fallback.
4. IF el AWS_Config_File no existe o no es legible, THEN THE Sincronizador_EKS SHALL mostrar un mensaje de error descriptivo y terminar con código de salida 1.
5. IF no se encuentran perfiles válidos en el AWS_Config_File, THEN THE Sincronizador_EKS SHALL mostrar un mensaje informativo indicando que no se encontraron perfiles y terminar con código de salida 0.

### Requisito 2: Descubrimiento de clústeres EKS

**User Story:** Como operador de Kubernetes, quiero que KSW descubra automáticamente los clústeres EKS en cada perfil AWS, para no tener que buscarlos manualmente en la consola de AWS.

#### Criterios de Aceptación

1. WHEN el Sincronizador_EKS procesa un Perfil_AWS, THE Sincronizador_EKS SHALL ejecutar el equivalente a `aws eks list-clusters --profile <nombre> --region <región>` para obtener la lista de clústeres EKS.
2. WHEN un Perfil_AWS retorna clústeres EKS, THE Sincronizador_EKS SHALL mostrar en la terminal el nombre del perfil, la región y la cantidad de clústeres encontrados.
3. IF las credenciales de un Perfil_AWS son inválidas o han expirado, THEN THE Sincronizador_EKS SHALL mostrar una advertencia para ese perfil y continuar con el siguiente perfil sin interrumpir la ejecución.
4. IF la llamada a la API de AWS falla por un error de red o timeout, THEN THE Sincronizador_EKS SHALL mostrar una advertencia para ese perfil y continuar con el siguiente perfil.
5. WHEN el Sincronizador_EKS completa el descubrimiento, THE Sincronizador_EKS SHALL mostrar un resumen con el total de clústeres encontrados en todos los perfiles.

### Requisito 3: Detección de duplicados en kubeconfig

**User Story:** Como operador de Kubernetes, quiero que KSW detecte los clústeres que ya existen en mi kubeconfig, para evitar entradas duplicadas y mantener mi configuración limpia.

#### Criterios de Aceptación

1. WHEN el Sincronizador_EKS encuentra un Clúster_EKS, THE Sincronizador_EKS SHALL verificar si ya existe un Contexto_Kubeconfig correspondiente en el Kubeconfig actual.
2. WHEN un Clúster_EKS ya tiene un Contexto_Kubeconfig existente, THE Sincronizador_EKS SHALL omitir ese clúster y mostrar un indicador de que fue saltado.
3. WHEN un Clúster_EKS no tiene un Contexto_Kubeconfig existente, THE Sincronizador_EKS SHALL marcarlo como pendiente de agregar.
4. THE Sincronizador_EKS SHALL comparar los clústeres usando el ARN del clúster EKS presente en el Contexto_Kubeconfig para determinar equivalencia.

### Requisito 4: Sincronización del kubeconfig

**User Story:** Como operador de Kubernetes, quiero que KSW agregue automáticamente los clústeres EKS faltantes a mi kubeconfig, para tener acceso inmediato a todos mis clústeres sin ejecutar comandos manualmente.

#### Criterios de Aceptación

1. WHEN el Sincronizador_EKS identifica un Clúster_EKS pendiente de agregar, THE Sincronizador_EKS SHALL ejecutar `aws eks update-kubeconfig --name <cluster> --profile <perfil> --region <región>` para agregarlo al Kubeconfig.
2. WHEN `aws eks update-kubeconfig` se ejecuta exitosamente, THE Sincronizador_EKS SHALL mostrar un indicador de éxito con el nombre del clúster y el perfil utilizado.
3. IF `aws eks update-kubeconfig` falla para un clúster específico, THEN THE Sincronizador_EKS SHALL mostrar un mensaje de error para ese clúster y continuar con los clústeres restantes.
4. WHEN la sincronización completa, THE Sincronizador_EKS SHALL mostrar un resumen final indicando: clústeres nuevos agregados, clústeres omitidos por duplicado y clústeres con error.

### Requisito 5: Validación de dependencias

**User Story:** Como operador de Kubernetes, quiero que KSW valide que las herramientas necesarias están instaladas antes de ejecutar la sincronización, para obtener un mensaje claro si falta algo.

#### Criterios de Aceptación

1. WHEN el usuario ejecuta `ksw eks kubeconfig`, THE Sincronizador_EKS SHALL verificar que el comando `aws` está disponible en el PATH del sistema antes de iniciar cualquier operación.
2. IF el comando `aws` no está disponible en el PATH, THEN THE Sincronizador_EKS SHALL mostrar un mensaje de error indicando que AWS CLI es requerido e incluir instrucciones de instalación, y terminar con código de salida 1.

### Requisito 6: Integración con la estructura de subcomandos de KSW

**User Story:** Como usuario de KSW, quiero que el nuevo comando siga la misma estructura y estilo visual que los demás subcomandos, para tener una experiencia consistente.

#### Criterios de Aceptación

1. THE KSW SHALL registrar `eks` como subcomando de primer nivel en el enrutador principal de `main()`.
2. WHEN el usuario ejecuta `ksw eks kubeconfig`, THE KSW SHALL delegar la ejecución al Sincronizador_EKS.
3. WHEN el usuario ejecuta `ksw eks` sin subcomando, THE KSW SHALL mostrar la ayuda de uso del subcomando `eks` con los subcomandos disponibles.
4. THE Sincronizador_EKS SHALL utilizar los mismos estilos de lipgloss (successStyle, warnStyle, dimStyle) definidos en KSW para mantener consistencia visual.
5. THE KSW SHALL incluir el subcomando `ksw eks kubeconfig` en la salida de `ksw -h`.

### Requisito 7: Soporte para filtrado por perfil

**User Story:** Como operador de Kubernetes, quiero poder sincronizar solo un perfil AWS específico, para ahorrar tiempo cuando solo necesito actualizar una cuenta.

#### Criterios de Aceptación

1. WHEN el usuario ejecuta `ksw eks kubeconfig --profile <nombre>`, THE Sincronizador_EKS SHALL procesar únicamente el Perfil_AWS especificado.
2. IF el perfil especificado con `--profile` no existe en el AWS_Config_File, THEN THE Sincronizador_EKS SHALL mostrar un mensaje de error indicando que el perfil no fue encontrado y terminar con código de salida 1.
3. WHEN el usuario ejecuta `ksw eks kubeconfig` sin la opción `--profile`, THE Sincronizador_EKS SHALL procesar todos los perfiles disponibles en el AWS_Config_File.
