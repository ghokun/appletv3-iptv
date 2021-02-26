# Apple TV 3 IPTV
This is an IPTV application for Apple TV 3 devices. It replaces RedbullTV app.

## Installation
1. Create DNS record for `appletv.redbull.tv` in your network.
2. Generate certificates for `appletv.redbull.tv` using script in scripts folder.
3. Download binary for your platform from releases.
4. Modify settings file and run
5. Install profile on Apple TV
6. Open RedbullTV application


## Credits
Code parts or ideas are taken from following repositories:
- https://github.com/iBaa/PlexConnect
- https://github.com/wahlmanj/sample-aTV
- https://github.com/jamesnetherton/m3u

## TASKS
- [x] Certificate generation
- [x] Intercept redbulltv
- [x] Profile installation
- [x] Parse m3u files
- [x] Play m3u8 files
- [x] Category and Channel images
- [x] Cache assets
- [x] Logger
- [x] Refresh m3u file from UI
- [x] Recently Played, 10 items
- [x] Search channels
- [x] Better templates
- [x] Localization
- [x] Accesibility labels
- [x] Cleanup and document go files
- [ ] Generate README
- [ ] Generate images for missing ones
- [ ] GitHub actions to autorelease
- [ ] EPG support
- [ ] Inject application icon
- [ ] Placeholder icons
- [ ] Cleanup javascript files
- [ ] Save favorites and recents to file
- [x] Embed files
- [ ] Prevent Apple TV software update
- [ ] Include DNS server