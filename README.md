## phi-service

Phigros score lookup service implemented using Go and Vue. Special thanks to the following projects:

- `Catrong/phi-plugin`
- `3035936740/Phigros_Resource`

When developing this tool, I referenced the code from these two contributors and made some modifications based on their work. Thank you for their open-source spirit OwO

## Usage

Make sure you have Golang and npm environment installed on your system, and place the latest version of Phigros installation package in `./build` directory (create it manually if necessary). Then execute the following commands:

```bash
make init && make && make run
```

Then the program will exit, and you need to fill in the configuration information. The configuration file is located at `./build/config.yaml`, and make sure to modify the paths of two files in the Data section to `./dist/xxx.csv`.

Finally, execute `make run` again. Then you can access the service through `http://localhost:[your_port]`.

## LICENSE

GNU General Public License v3.0

This project is open source under the GPLv3 license. You are free to use, modify, and distribute this project under the condition of complying with the license.
