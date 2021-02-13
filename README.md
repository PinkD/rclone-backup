# Backup

backup your directories with rclone

## Usage

- setup rclone
    - you need to download [rclone](https://github.com/rclone/rclone) and set up a remote
- add config
    - see `config` section for detail
- just run it

## Config

template:

```json
[
  {
    "name": "-",
    "remote": "gd:/backup",
    "interval": "1h",
    "proxy": "http://127.0.0.1:1080"
  },
  {
    "name": "path1",
    "path": "/path1"
  },
  {
    "name": "path2",
    "path": "/path2",
    "interval": "1d",
    "remote": "gd:/project",
    "proxy": "https://127.0.0.1:1080"
  }
]
```

- `name` 
  - `-` for global config
    - for example, `remote` in `path1` will be `gd:/backup`
- `remote`
  - your rclone remote
- `interval`
  - run interval, if a work is still running, it won't be triggered again
- `proxy` (optional)
  - proxy for rclone, only support http and https (**no socks5**)

## License

```license
   Copyright 2021 PinkD

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
```
