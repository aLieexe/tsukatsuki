## Using one server for multiple project
Tsukatsuki can work with two project, as long as the proxy share the same type, when initializing just make sure you pick the same proxy for the two project, then proceed just how you would deploy the first project.


## Problems that you may encounter
1. WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!

    simply follow the command that is given by that error and it should solve it, it looks like this
    `ssh-keygen -f '/home/user/.ssh/known_hosts' -R 'ip-address'`


2. Network is unreachable

    This happen because no matching route or default gateway for the destination network/address family (interface down, missing IPv6)
    You can test for IPv6 here <https://test-ipv6.com/> . If you dont have one, the easiest way to get is by using Cloudflare Warp

3. No route to host

    A route exists but the specific host is unreachable, this could happen the original SSH port is already changed / blocked by UFW
    You will need to re-specify the correct one in `tsukatsuki.yaml`

