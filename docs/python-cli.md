# Member 增删改查

计划使用 python 来写一个交互式的简易脚本
在服务器上直接对成员进行增删改查

golang的服务端则采用 双 router 模式
其中 internal router 只监听服务器的端口
不通过反向代理进行暴露

---

## 添加成员

`POST` `/member/create`

```json
{
    "power": 1,
    "nick": "nick_name",
    "email": "email_str",
    "password": "password_original",
    "is_delete": false
}
```

---

## 更新

`PATCH` `/member/update`

```json
{
    "id": 1,
    "nick": "nick_name",
    "email": "email_str",
    "password": "password_original",
    "is_delete": false
}
```
