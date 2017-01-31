# emptynest
An American sitcom that originally aired on NBC from October 8, 1988 to April 29, 1995.


### Configuration
Example
```
server_addr = ":8000"
server_debug = true
db_file = "data.db"
log_dir = "./logs"
data_dir = "./data"
get_location = "query"
get_param = "JSESSIONID"
post_location = "data"
post_param="data"
host_info_plugin = "info.so"
payload_plugin_directories = ["./plugins"]
encoder_plugin_chain = ["zip.so", "base64.so"]
crypto_plugin_chain = ["rc4.so", "aes.so"]
key_chain = ["AAAA", "lfRH8Vp90iqHj2YPR0Kdw3Xi423AEcv6"]
kill_date = ""
```
