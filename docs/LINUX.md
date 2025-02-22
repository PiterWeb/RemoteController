# Linux

## Execute in Linux

When you run RemoteController you will need to be sure to:

- Your user has read/write permissions for /dev/input/event/* and uinput devices
    - Example in Debian: `sudo usermod -aG input $USER`
- Uinput module enabled
    - Example in Debian: `sudo modprobe uinput`
