# ksw aliases demo
Output demo/aliases.gif

Set Shell "bash"
Set FontSize 14
Set Width 800
Set Height 300
Set Theme "Dracula"
Set Padding 20
Set TypingSpeed 80ms
Set PlaybackSpeed 0.8

Env KUBECONFIG demo/kubeconfig-demo.yaml

Type "ksw alias prod eks-payments-pdn"
Sleep 400ms
Enter
Sleep 1000ms

Type "ksw @prod"
Sleep 400ms
Enter
Sleep 1200ms

Type "ksw alias ls"
Sleep 400ms
Enter
Sleep 1200ms
