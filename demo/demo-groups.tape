# ksw groups demo
Output demo/groups.gif

Set Shell "bash"
Set FontSize 14
Set Width 800
Set Height 360
Set Theme "Dracula"
Set Padding 20
Set TypingSpeed 80ms
Set PlaybackSpeed 0.8

Env KUBECONFIG demo/kubeconfig-demo.yaml

Type `ksw group add payments "eks-payments*"`
Sleep 400ms
Enter
Sleep 1200ms

Type "ksw group ls"
Sleep 400ms
Enter
Sleep 1200ms

Type "ksw group use payments"
Sleep 400ms
Enter
Sleep 1500ms

Escape
Sleep 800ms
