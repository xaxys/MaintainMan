# MaintainMan

[![OurEDA 2022](https://img.shields.io/badge/OurEDA-2022-00ffcc.svg)](https://img.shields.io/badge/OurEDA-2022-00ffcc)
[![License](https://img.shields.io/badge/license-MIT%20with%20PATENTS-green.svg)](https://github.com/xaxys/MaintainMan/blob/master/LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/xaxys/MaintainMan/pulls)
[![Contributors](https://img.shields.io/github/contributors/xaxys/maintainman.svg)](https://github.com/xaxys/MaintainMan/graphs/contributors)
[![Go Version](https://img.shields.io/github/go-mod/go-version/xaxys/maintainman.svg)](https://github.com/xaxys/MaintainMan/blob/master/go.mod)
[![Release](https://img.shields.io/github/v/release/xaxys/maintainman.svg)](https://github.com/xaxys/MaintainMan/releases)
[![Downloads](https://img.shields.io/github/downloads/xaxys/maintainman/total.svg)](https://github.com/xaxys/MaintainMan/releases)
[![Build and Test](https://github.com/xaxys/MaintainMan/actions/workflows/main.yml/badge.svg)](https://github.com/xaxys/MaintainMan/actions/workflows/main.yml)

MaintainMan is a logistic report management system powered by iris.

## Feature

- RESTful HTTP API

- User management with configurable Role-Based access control

- Database: Mysql, Sqlite3

- Storage: S3, Local

- Cache: Redis, Local

- 3 pulggable modules

  - Order management

    - 8 status available

      - Waiting for Assignment

      - Order Assigned

      - Order Completed

      - Order Appraised

      - Reported as pending

      - Hold for a while

      - Order Canceled

      - Order Rejected

    - Switchable order comment

    - Order assignment system

      - One Repairer for one assignment

      - Supports multiple order assignments

    - Order appraising system and performance display

    - Item inventory management associated with the order system

  - Announcement management
  
    - Configurable display times

    - User click statistics

  - Image Hosting service
  
    - Auto watermarking

    - Custom transformation

      - Resizing & Croping

      - Text with color and font

    - Image compression

    - Configurable image cache

- More...

## Configuration

MaintainMan has 6 configuration files now. All configuration files have independent version control.

When the maintainman detected a old version configuration file, it will automatically upgrade it (conflict field will be skipped).

When the maintainman detected a new version configuration file, it will send a warning message.

### app.yml

App config is used to configure the database and various connection parameters as well as the functional parameters of the core system.

<details>
<summary>example</summary>

```yaml
app:
  # application name.
  name: "maintainman"
  # listen port.
  listen: ":8080"
  # log level (debug, info, warn, error, fatal).
  loglevel: "debug"

  page:
    # max number of items in a page.
    # http request paramenter `limit` must not exceed this value.
    limit: 100
    # default number of items in a page.
    # this number will be used when http request paramenter `limit`
    # is not specified or <= 0.
    default: 50

token:
  # token secret.
  # IMPORTANT! you'd better change it to a random string or a strong
  # secret key.
  key: ""
  # token expire duration.
  expire: "30m"

database:
  # database type (mysql, sqlite).
  driver: "mysql"
  mysql:
    host: "localhost"
    port: 3306
    name: "maintainman"
    params: "parseTime=true&loc=Local&charset=utf8mb4"
    user: "root"
    password: ""
  sqlite:
    # sqlite database file path.
    path: "maintainman.db"

storage:
  # storage type (local, s3).
  driver: "local"
  local:
    path: "./files"
  # if s3 connection defined here, module config without s3 connection
  # will use the connection defined here.
  s3:
    access_key: ""
    secret_key: ""
    bucket: ""
    region: ""

cache:
  # cache type (local, redis).
  driver: "local"
  # cache limit. if the cache limit is reached, some entries will be
  # evicted automatically.
  # if the cache limit is 0, no entries will not be evicted.
  limit: 268435456 # 256M
  redis:
    host: "localhost"
    port: 6379
    password: ""

throttling:
  enable: false
  # the max number of requests allowed in a period.
  burst: 100
  # the duration between requests.
  rate: 10
  # the purge duration.
  purge: 1m
  # the expiration duration.
  expire: 10m

# enabled modules
module:
  role: true
  user: true
  image: true
  announce: true
  order: true
  wxnotify: true
  word: true
  sysinfo: true

# channel size of event bus (message bus).
bus_buffer: 1000

```

</details>

### user.yml

User config is used to configure login and user management options.

<details>
<summary>example</summary>

```yaml
wechat:
  # wechat appid.
  appid: ""
  # wechat secret.
  secret: ""
  # whether a unregistered user will be registered on wechat login.
  # if false, reponse code will be `403` when a unregistered user try
  # wechat login.
  # if true, a unregistered user will be registered on wechat login.
  # username will be open_id and user will be assigned a random password.
  fastlogin: true

cache:
  # cache type (local, redis).
  driver: local
  # cache limit.
  limit: 268435456 # 256M
  # if redis, connection has been configured in app.yml

# the admin user configuration.
# only apply at first initialization.
# IMPORTANT! you'd better change it to some strong password and delete
# belowing entries after the first initialization.
admin:
  name: "admin"
  display_name: "maintainman default admin"
  role_name: "super_admin"
  password: "12345678"

```

</details>

### role.yml

Role config is used to configure all roles and their corresponding permissions. Roles are ordered. Only buttom-up inheritance is valid (latter roles are superior).

<details>
<summary>example</summary>

```yaml
role:

- display_name: 封停用户
  name: banned
  permissions: []
  inheritance: []

- name: guest
  display_name: 访客
  guest: true
  permissions:
  - user.register
  - user.login
  - user.wxlogin
  - user.wxregister
  inheritance: []

- name: user
  display_name: 普通用户
  default: true
  permissions:
  - image.upload
  - image.view
  - user.view
  - user.update
  - user.renew
  - role.view
  - announce.view
  - announce.hit
  - order.view
  - order.create
  - order.cancel
  - order.update
  - order.appraise
  - order.urgence
  - order.comment.view
  - order.comment.create
  - order.comment.delete
  - tag.view.1
  - tag.add.1
  # `tag.add.1` is a special permission.
  # in `perm.?` pattern, if `?` is a number, the number will be compared
  # to judge whether the role has the permission.
  # e.g. if a role has `perm.2`, then the `perm.2` and `perm.1` will be
  # judge as true.
  inheritance:
  - guest

- name: maintainer
  display_name: 维护工
  permissions:
  - order.viewfix
  - order.reject
  - order.report
  - order.complete
  - item.consume
  - item.viewall
  - tag.view.2
  - tag.add.2
  inheritance:
  - user

- name: super_maintainer
  display_name: 维护工（可自行接单）
  permissions:
  - order.selfassign
  - order.viewall
  inheritance:
  - maintainer

- name: admin
  display_name: 管理员
  permissions:
  - image.*
  - division.*
  - announce.*
  - order.*
  - tag.*
  - item.*
  # in `perm.*` pattern, `*` means any, all sub permissions under perm will
  # be judged as true.
  inheritance:
  - maintainer

- name: super_admin
  display_name: 超级管理员
  permissions:
  - '*'
  inheritance:
  - admin

```

</details>

### image.yml

Image config is used to configure image hosting service and predefined transformations.

<details>
<summary>example</summary>

```yaml
# jpeg compression quality.
jpeg_quality: 80
# max gif color number.
gif_num_colors: 256
# all image after transformation will be cached as jpeg.
cache_as_jpeg: true
# all image uploaded will be saved as gif.
save_as_jpeg: false

upload:
  # upload request returns straight after image is processed by the server.
  # but saving might still fail.
  async: false
  # the max file size of image allowed to upload.
  max_file_size: 10485760 # 10 MB
  # the max dimension of image allowed to upload.
  max_pixels: 15000000    # 15 million pixels
  # the throttling rate control.
  throttling:
    enable: true
    # the max number of requests allowed in a period.
    burst: 20
    # the duration between requests.
    rate: 5
    # the purge duration.
    purge: 1m
    # the expiration duration.
    expire: 5m

cache:
  # cache type (local, redis).
  driver: local
  # cache limit. if the cache limit is reached, image in storage
  # will be deteted automatically.
  # if the cache limit is 0, no entries will not be evicted.
  # (strongly not recommended)
  limit: 1073741824 # 1 GB
  # if redis, connection has been configured in app.yml

storage:
  # storage type (local, s3).
  driver: local
  local:
    path: ./images
  s3:
    # if access_key and secret_key are not set, s3 connection defined
    # in app.yml will be used.
    # access_key: ""
    # secret_key: ""
    # region: ""
    bucket: "Image"
  # image cache storage. sub path of main storage.
  # e.g. if main storage is ./images, cache storage is ./images/cache,
  cache:
    # whether the storage path will be cleaned up on server start.
    # recommended to be true if you are using local cache instead of redis.
    clean: true

transformations:
  # predefined transformations.
  # square returns a 256 x 256 square image chopped from the center.
  square:
    params:   w_256,h_256,c_p,g_c
    # Run on every upload
    eager:   true
  # watermarked returns a equal scaling, 800 widthm, watermarked image.
  watermarked:
    # if params is not set, the transformation will be applied on.
    default: true
    params: w_800
    texts:
    # text will be added to the bottom right corner of the image.
    # the {{.Name}} will be replaced by the upload user name.
    - content: "{{.Name}}@MaintainMan"
      gravity: se
      # text position in the image. relative to gravity.
      # non-negative integer.
      x-pos:   10
      y-pos:   0
      # color format is hex.
      # e.g. #RRGGBBAA or #RRGGBB or #RGBA or #RGB
      color:   "#808080CC"
      # font file path. if not set, will search filename in
      # embedded fonts.
      font:    fonts/SourceHanSans-Regular.ttf
      size:    14

```

</details>

### announce.yml

Announce config is used to configure announcement management.

<details>
<summary>example</summary>

```yaml
# the duration that a user hit the same announcement will not
# be counted repeatedly.
hit_expire: "12h"

cache:
  driver: "local"
  limit: 268435456 # 256M

```

</details>

### order.yml

Order config is used to configure order management.

<details>
<summary>example</summary>

```yaml
# Whether item count can be negative.
# if false, an `item count is not enough` error may be returned on
# item consuming.
item_can_negative: true

appraise:
  # the duration that a user can appraise the order after the
  # order completed.
  # the order will be appraised automatically after the duration.
  timeout: "72h"
  # the duration that the system will check the timeouted unappraised order.
  purge: "10m"
  # the default appraise score of timeouted unappraised order.
  default: 5

notify:
  wechat:
    status:
      tmpl:   "微信订阅消息模板id"
      order:  "模板中 订单编号 字段名"
      title:  "模板中 订单标题 字段名"
      status: "模板中 订单状态 字段名"
      time:   "模板中 订单更新时间 字段名"
      other:  "模板中 备注 字段名 (用于传递维修工信息)"

    comment:
      tmpl:    "微信留言消息模板id"
      title:   "模板中 订单标题 字段名"
      name:    "模板中 留言人 字段名"
      messgae: "模板中 留言内容 字段名"
      time:    "模板中 留言时间 字段名"

```

</details>

## Documentation

Find document here [Maintainman Doc](https://maintainman.oasis.run/).

Or On [Github Wiki](https://github.com/xaxys/MaintainMan/wiki/API-Docs).

## Contributing

You can help to make the project better by creating an issue or pull request.

## Author

- xaxys

  - [github.com/xaxys](https://github.com/xaxys)

  - [alappm@qq.com](mailto:alappm@qq.com)

- Marksagittarius

  - [github.com/Marksagittarius](https://github.com/Marksagittarius)

- DawningW

  - [github.com/DawningW](https://github.com/DawningW)

## License

Released under [MIT License](https://github.com/xaxys/MaintainMan/blob/master/LICENSE). We also provide an additional [patent grant](https://github.com/xaxys/MaintainMan/blob/master/PATENTS).
