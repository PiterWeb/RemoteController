# Linux

## Execute in Linux

When you run RemoteController you will need to be sure to:

- Your user has read/write permissions for /dev/input/event/* and uinput devices
    - Example in Debian:
      ```sh
      sudo usermod -aG input $USER
      ```
- Uinput module enabled
    - Check if it is loaded:
      - Example in Debian:
          ```sh
          lsmod | grep uinput
          ```
    - Load the module:
      - Example in Debian:
          ```sh
          sudo modprobe uinput
          ```
