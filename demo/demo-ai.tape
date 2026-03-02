# ksw AI demo — natural language switching
Output demo/ai.gif

Set Shell "bash"
Set FontSize 14
Set Width 800
Set Height 420
Set Theme "Dracula"
Set Padding 20
Set TypingSpeed 80ms
Set PlaybackSpeed 0.8

Env KUBECONFIG demo/kubeconfig-demo.yaml

Type `ksw ai "switch to payments pdn"`
Sleep 400ms
Enter
Sleep 2000ms

Type `ksw ai "create a group called backend with payments and orders"`
Sleep 400ms
Enter
Sleep 2500ms

Type `ksw ai "go back to the previous one"`
Sleep 400ms
Enter
Sleep 2000ms
