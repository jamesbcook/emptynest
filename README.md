# emptynest
Emptynest is a plugin based C2 server framework. The goal of this project is not to replace robust tools such as Empire, Metasploit, or Cobalt Strike. Instead, the goal is to create a supporting framework for quickly creating small, purpose built handlers for custom agents. No agent is provided. Users of Emptynest should create their own agents that implement minimal functionality and can be used to evade detection and establish a more robust channel. An example of an agent might support Unhooking, DLL Unloading, and code execution. Due to the simple nature of this project, it is recommended that agents be kept private.

This project was originally created to bypass sandboxed execution environments, and borrows ideas and research from [Ebowla](https://github.com/Genetic-Malware/Ebowla). Instead of using keyed payloads, handlers use an approve/deny process to allow an operator to identify and prevent continued execution in a sandboxed environment. For example, an info plugin can be used that sends a username and hostname for a Windows target. The operator could see a request come in that looks like the following:
```
[!] APPROVAL REQUESTED:
ID: 1
Info:
Username: TENFORWARD\tom
Hostname: tenforward
```

And another:
```
[!] APPROVAL REQUESTED:
Info:
Username: JOHN-PC\Admin
Hostname: axjekaufjaksd
```

Depending on the level of sophistication of the sandbox and information known about target systems, this could be easy to discern between the two and approve one while denying the other. Info plugins can be easily adapted to send additional details; including Active Directory, running processes, memory usage, etc.

Another goal when developing this project was to enable chainable encryption and encoding schemes that can be changed quickly. An idea borrowed from Russian APT29 as discussed in [No Easy Breach - DerbyCon 2016](https://www.slideshare.net/MatthewDunwoody1/no-easy-breach-derby-con-2016). An example chain could be RC4->AES-CTR->BASE64. The server configuration would look like:
```
encoder_plugin_locations = ["base64.so"] # encoding plugins to use
crypto_plugin_locations = ["rc4.so", "aes_ctr.so"] # encryption plugins to use
key_chain = ["41414141", "6c66524838567039306971486a32595052304b64773358693432334145637636"] # encryption keys in order of plugin
```
The possibilities are endless. For example, you could modify RC4 plugin to brute-force an incoming key. You may notice the server plugins provided in the repo do not implemented authenticated encryption, a principal you can read more about [here](https://moxie.org/blog/the-cryptographic-doom-principle/). We did this to keep message lengths to a minimum, again, the idea behind this project was to have minimal functionality and minimal traffic between the server and agent. You can easily modify the encryption plugins or handler itself to validate the message integrity should you feel that is needed.
