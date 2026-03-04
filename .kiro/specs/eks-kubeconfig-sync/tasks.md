# Plan de Implementación: eks-kubeconfig-sync

## Resumen

Implementar el subcomando `ksw eks kubeconfig` en un nuevo archivo `eks.go`. El desarrollo es incremental: primero las estructuras y funciones puras (parseables y testeables), luego la lógica de descubrimiento y sincronización, y finalmente la integración con `main.go`. Los tests de propiedades usan `pgregory.net/rapid`.

## Tareas

- [x] 1. Crear archivo `eks.go` con estructuras de datos y funciones de parseo
  - [x] 1.1 Crear `eks.go` con las estructuras `awsProfile`, `eksCluster` y `syncResult`
    - Definir las tres structs según el diseño
    - Incluir imports necesarios (`bufio`, `encoding/json`, `fmt`, `os`, `os/exec`, `strings`)
    - _Requisitos: 1.1, 2.1, 4.4_

  - [x] 1.2 Implementar `parseAWSProfiles(configPath string) ([]awsProfile, error)` para leer perfiles de `~/.aws/config`
    - Parsear secciones `[default]` y `[profile <nombre>]` con `bufio.Scanner`
    - Extraer la clave `region` de cada perfil
    - Retornar error si el archivo no existe o no es legible
    - Retornar lista vacía si no hay perfiles válidos
    - _Requisitos: 1.1, 1.2, 1.4, 1.5_

  - [x] 1.3 Implementar `getDefaultRegion() string` para obtener la región fallback
    - Ejecutar `aws configure get region` y retornar el resultado
    - Si falla, retornar `"us-east-1"` como fallback
    - _Requisito: 1.3_

  - [x] 1.4 Implementar `filterProfiles(profiles []awsProfile, profileFilter string) ([]awsProfile, error)` para filtrado por `--profile`
    - Si `profileFilter` está vacío, retornar todos los perfiles
    - Si tiene valor, retornar solo el perfil que coincida o error si no existe
    - _Requisitos: 7.1, 7.2, 7.3_

  - [ ] 1.5 Escribir test de propiedad para parseo de perfiles AWS
    - **Propiedad 1: Round-trip de parseo de perfiles AWS**
    - Generar archivos AWS config aleatorios con N perfiles, parsear, verificar que se extraen exactamente N perfiles con nombres y regiones correctas
    - Usar `pgregory.net/rapid` en archivo `eks_test.go`
    - **Valida: Requisitos 1.1, 1.2**

  - [ ] 1.6 Escribir test de propiedad para filtrado de perfiles
    - **Propiedad 6: Filtrado de perfiles por nombre**
    - Generar listas de perfiles y filtros, verificar que filtro vacío retorna todos y filtro específico retorna solo el perfil correcto
    - **Valida: Requisitos 7.1, 7.3**

- [x] 2. Implementar descubrimiento de clústeres EKS y detección de duplicados
  - [x] 2.1 Implementar `checkAWSCLI() error` para validar que `aws` está en el PATH
    - Usar `exec.LookPath("aws")` para verificar disponibilidad
    - Retornar error con mensaje descriptivo e instrucciones de instalación si no está disponible
    - _Requisitos: 5.1, 5.2_

  - [x] 2.2 Implementar `listEKSClusters(profile, region string) ([]string, error)` para descubrir clústeres
    - Ejecutar `aws eks list-clusters --profile <profile> --region <region> --output json`
    - Parsear el JSON de respuesta para extraer la lista de nombres de clústeres
    - Manejar errores de credenciales y red sin interrumpir la ejecución
    - _Requisitos: 2.1, 2.3, 2.4_

  - [x] 2.3 Implementar `getExistingEKSContexts() (map[string]bool, error)` para leer contextos existentes del kubeconfig
    - Ejecutar `kubectl config get-contexts -o name` para obtener los contextos actuales
    - Filtrar contextos que contengan `arn:aws:eks:` y almacenarlos en un mapa
    - _Requisitos: 3.1, 3.4_

  - [x] 2.4 Implementar `buildClusterARN(cluster, region, accountID string) string` para construir ARNs
    - Construir el ARN con formato `arn:aws:eks:<region>:<account>:cluster/<name>`
    - _Requisito: 3.4_

  - [x] 2.5 Implementar función de parseo de JSON de `list-clusters` como función pura testeable
    - Extraer la función de parseo del JSON para que sea testeable independientemente de `listEKSClusters`
    - _Requisito: 2.1_

  - [ ]* 2.6 Escribir test de propiedad para parseo de output de list-clusters
    - **Propiedad 2: Parseo de output de list-clusters**
    - Generar listas JSON de nombres de clústeres, parsear, verificar que se retornan exactamente los mismos nombres en el mismo orden
    - **Valida: Requisito 2.1**

  - [ ]* 2.7 Escribir test de propiedad para construcción y matching de ARN
    - **Propiedad 4: Consistencia de construcción y matching de ARN**
    - Generar combinaciones de nombre/región/account, construir ARN, verificar formato `arn:aws:eks:<region>:<account>:cluster/<name>` y que la función de matching lo reconozca
    - **Valida: Requisito 3.4**

- [x] 3. Checkpoint - Verificar que las funciones puras y tests pasan
  - Ejecutar `go test ./...` y asegurar que todos los tests pasan. Preguntar al usuario si surgen dudas.

- [x] 4. Implementar lógica de sincronización y resumen
  - [x] 4.1 Implementar `updateKubeconfig(cluster, profile, region string) error` para agregar clústeres al kubeconfig
    - Ejecutar `aws eks update-kubeconfig --name <cluster> --profile <profile> --region <region>`
    - Retornar error si el comando falla
    - _Requisitos: 4.1, 4.2, 4.3_

  - [x] 4.2 Implementar la lógica de partición de clústeres nuevos vs existentes
    - Comparar los clústeres descubiertos contra los contextos existentes usando ARNs
    - Clasificar cada clúster como "nuevo" o "existente"
    - _Requisitos: 3.1, 3.2, 3.3_

  - [ ]* 4.3 Escribir test de propiedad para partición de clústeres
    - **Propiedad 3: Partición correcta de clústeres nuevos vs existentes**
    - Generar conjuntos de clústeres descubiertos y contextos existentes, verificar que la partición es disjunta y su unión es igual al conjunto original
    - **Valida: Requisitos 3.1, 3.2, 3.3**

  - [ ]* 4.4 Escribir test de propiedad para invariante del resumen
    - **Propiedad 5: Invariante del resumen de sincronización**
    - Generar tripletas (added, skipped, failed), verificar que `added + skipped + failed == total`
    - **Valida: Requisito 4.4**

- [x] 5. Implementar `handleEksKubeconfig` y `handleEks` como orquestadores
  - [x] 5.1 Implementar `handleEksKubeconfig(profileFilter string)` como flujo principal
    - Llamar a `checkAWSCLI()`, si falla mostrar error con `warnStyle` y salir con código 1
    - Llamar a `parseAWSProfiles()`, manejar errores y caso de lista vacía
    - Aplicar `filterProfiles()` si se proporcionó `--profile`
    - Para cada perfil: llamar a `listEKSClusters()`, mostrar progreso con estilos de lipgloss
    - Obtener contextos existentes con `getExistingEKSContexts()`
    - Para cada clúster nuevo: llamar a `updateKubeconfig()`, mostrar resultado con `successStyle` o `warnStyle`
    - Para clústeres existentes: mostrar con `dimStyle` que fueron omitidos
    - Mostrar resumen final con contadores de `syncResult`
    - _Requisitos: 1.1, 1.4, 1.5, 2.2, 2.3, 2.4, 2.5, 3.2, 3.3, 4.1, 4.2, 4.3, 4.4, 5.1, 5.2, 6.4_

  - [x] 5.2 Implementar `handleEks()` como enrutador de subcomandos de `ksw eks`
    - Parsear `os.Args` para detectar subcomando `kubeconfig` y flag `--profile`
    - Si no hay subcomando, mostrar ayuda de uso del subcomando `eks`
    - Delegar a `handleEksKubeconfig(profileFilter)` cuando corresponda
    - _Requisitos: 6.2, 6.3_

- [x] 6. Integrar con `main.go` y actualizar ayuda
  - [x] 6.1 Agregar case `"eks"` en el switch de `main()` en `main.go`
    - Agregar `case "eks": handleEks(); return` en el switch de `os.Args[1]`
    - _Requisito: 6.1_

  - [x] 6.2 Actualizar el bloque de ayuda `-h` en `main.go`
    - Agregar las líneas de ayuda para `ksw eks kubeconfig` y `ksw eks kubeconfig --profile <name>`
    - _Requisito: 6.5_

- [x] 7. Checkpoint final - Verificar integración completa
  - Ejecutar `go build -o ksw .` para verificar que compila sin errores. Ejecutar `go test ./...` para verificar que todos los tests pasan. Preguntar al usuario si surgen dudas.

## Notas

- Las tareas marcadas con `*` son opcionales y pueden omitirse para un MVP más rápido
- Cada tarea referencia los requisitos específicos para trazabilidad
- Los checkpoints aseguran validación incremental
- Los tests de propiedades usan `pgregory.net/rapid` y validan invariantes universales
- Los tests unitarios validan ejemplos concretos y edge cases
- Todo el código nuevo va en `eks.go` y `eks_test.go`, solo se modifica `main.go` para integración
