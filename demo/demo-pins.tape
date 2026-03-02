# ksw pins demo
Output demo/pins.gif

Set Shell "bash"
Set FontSize 14
Set Width 800
Set Height 300
Set Theme "Dracula"
Set Padding 20
Set TypingSpeed 80ms
Set PlaybackSpeed 0.8

Env KUBECONFIG demo/kubeconfig-demo.yaml

Type "ksw pin payments-dev"
Sleep 400ms
Enter
Sleep 1000ms

Type "ksw pin orders-dev"
Sleep 400ms
Enter
Sleep 1000ms

Type "ksw pin ls"
Sleep 400ms
Enter
Sleep 1200ms

Type "ksw pin rm orders-dev"
Sleep 400ms
Enter
Sleep 1000ms
