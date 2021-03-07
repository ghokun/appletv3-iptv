# Apple TV 3 IPTV
This is an IPTV application for Apple TV 3 devices. It replaces RedbullTV app.

## Installation
1. Create DNS record for `appletv.redbull.tv` in your network.
> `appletv.redbull.tv` should point to ip address that `appletv3-iptv` runs.
2. Generate certificates for `appletv.redbull.tv`
```bash
openssl req -new -nodes -newkey rsa:2048 -out redbulltv.pem -keyout redbulltv.key -x509 -days 7300 -subj "/C=US/CN=appletv.redbull.tv"
openssl x509 -in redbulltv.pem -outform der -out redbulltv.cer && cat redbulltv.key >> redbulltv.pem
```
3. Download binary for your platform from releases.
4. Create a settings file and run
```yaml
# See sample/config.yaml
---
# You can leave m3u link empty and set it from settings in app
m3uPath: ./sample/sample.m3u # or https://domain.com/sample.m3u
httpPort: "80"
httpsPort: "443"
cerPath: ./sample/certs/redbulltv.cer
pemPath: ./sample/certs/redbulltv.pem
keyPath: ./sample/certs/redbulltv.key
logToFile: true
loggingPath: log
recents: []
favorites: []
```
```bash
./appletv3-iptv -config config.yaml
```
5. Install profile on Apple TV
```
1. Open Apple TV
2. Go to Settings > General
3. Set Send Data to Apple to `No`.
4. Press `Play` button on Send Data to Apple
5. Add Profile > Ok
6. Enter URL: http://appletv.redbull.tv/redbulltv.cer
```
6. Open RedbullTV application


## Credits
Code parts or ideas are taken from following repositories:
- https://github.com/iBaa/PlexConnect
- https://github.com/wahlmanj/sample-aTV
- https://github.com/jamesnetherton/m3u

## Tasks
- [ ] Cleanup javascript files
- [ ] Inject application icon
- [ ] EPG support
- [ ] Include DNS server
- [ ] Prevent Apple TV software update