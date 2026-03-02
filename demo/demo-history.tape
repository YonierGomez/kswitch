# ksw history demo
Output demo/history.gif

Set Shell "bash"
Set FontSize 14
Set Width 800
Set Height 300
Set Theme "Dracula"
Set Padding 20
Set TypingSpeed 80ms
Set PlaybackSpeed 0.8

Env KUBECONFIG demo/kubeconfig-demo.yaml

Type "ksw payments-dev"
Sleep 400ms
Enter
Sleep 800ms

Type "ksw orders-qa"
Sleep 400ms
Enter
Sleep 800ms

Type "ksw billing-dev"
Sleep 400ms
Enter
Sleep 800ms

Type "ksw history"
Sleep 400ms
Enter
Sleep 1200ms

Type "ksw history 2"
Sleep 400ms
Enter
Sleep 1000ms
