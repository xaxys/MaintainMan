# MaintainMan

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
version: 1.2.2
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
    path: "./images"
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

### permission.yml

Permission config is used to configure all permissions and their corresponding names. It only affects the permission name display.
No permission will be added or deleted if you change the config file.

<details>
<summary>example</summary>

```yaml
version: 1.2.0
permission:
  announce:
    create: 创建公告
    delete: 删除公告
    hit: 点击公告
    update: 更新公告
    view: 查看公告
    viewall: 查看所有公告
  comment:
    create: 创建评论
    createall: 创建所有评论
    delete: 删除评论
    deleteall: 删除所有评论
    view: 查看我的评论
    viewall: 查看所有评论
  division:
    create: 创建分组
    delete: 删除分组
    update: 更新分组
    viewall: 查看所有分组
  image:
    custom: 处理图片
    upload: 上传图片
    view: 查看图片
  item:
    consume: 消耗零件
    create: 创建零件
    delete: 删除零件
    update: 更新零件
    viewall: 查看所有零件
  order:
    appraise: 评分
    assign: 分配订单
    cancel: 取消订单
    comment:
      create: 创建评论
      createall: 创建所有评论
      delete: 删除评论
      deleteall: 删除所有评论
      view: 查看我的评论
      viewall: 查看所有评论
    complete: 完成订单
    create: 创建订单
    defect: 修改故障分类
    hold: 挂起订单
    reject: 拒绝订单
    release: 释放订单
    report: 上报订单
    selfassign: 给自己分配订单
    update: 更新订单
    updateall: 更新所有订单
    urgence: 修改紧急程度
    view: 查看我的订单
    viewall: 查看所有订单
    viewfix: 查看我维修的订单
  permission:
    viewall: 查看所有权限
  role:
    create: 创建角色
    delete: 删除角色
    update: 更新角色
    view: 查看当前角色
    viewall: 查看所有角色
  tag:
    add: 添加标签
    create: 创建标签
    delete: 删除标签
    view: 查看标签
  user:
    create: 创建用户
    delete: 删除用户
    division: 修改部门
    login: 登录
    register: 注册
    renew: 更新Token
    role: 修改角色
    update: 更新用户
    updateall: 更新所有用户
    view: 查看当前用户
    viewall: 查看所有用户
    wxlogin: 微信登录
    wxregister: 微信注册

```

</details>

### role.yml

Role config is used to configure all roles and their corresponding permissions. Roles are ordered. Only buttom-up inheritance is valid (latter roles are superior).

<details>
<summary>example</summary>

```yaml
version: 1.2.0

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
    # the max number of requests allowed in a period.
    burst: 20
    # the duration between requests.
    rate: 1
    # the purge duration.
    purge: 1m
    # the expiration duration.
    expire: 1m

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

Released under [**Modified** MIT License](https://github.com/xaxys/MaintainMan/blob/master/LICENSE).
