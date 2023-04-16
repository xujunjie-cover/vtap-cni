

## FAQ
0. Why it appear?
A: As we know, Qemu use Tap device as a nic. When use some projects running on k8s like Kubevirt, kata, we have to create a tap device in netns. Kubevirt create a extra bridge. a dummy nic, a tap nic; 



