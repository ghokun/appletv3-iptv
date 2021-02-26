# Apple TV 3 IPTV
This is an IPTV application for Apple TV 3 devices. It replaces RedbullTV app.

## Installation
1. Create DNS record for `appletv.redbull.tv` in your network.
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
m3uPath: sample.m3u # or https://domain.com/sample.m3u
httpPort: "80"
httpsPort: "443"
cerPath: certs/redbulltv.cer
pemPath: certs/redbulltv.pem
keyPath: certs/redbulltv.key
logToFile: true
loggingPath: log
favorites: []
```
```bash
./appletv3-iptv -config config.yaml
```
5. Install profile on Apple TV
6. Open RedbullTV application


## Credits
Code parts or ideas are taken from following repositories:
- https://github.com/iBaa/PlexConnect
- https://github.com/wahlmanj/sample-aTV
- https://github.com/jamesnetherton/m3u

## Tasks
- [ ] GitHub actions to autorelease
- [ ] Placeholder icons
- [ ] Cleanup javascript files
- [ ] Save favorites and recents to file
- [ ] Include DNS server
- [ ] Prevent Apple TV software update
- [ ] Inject application icon
- [ ] EPG support