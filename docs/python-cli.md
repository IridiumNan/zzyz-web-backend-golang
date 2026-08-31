# Member 增删改查

计划使用 python 来写一个交互式的简易脚本
在服务器上直接对成员进行增删改查

golang的服务端则采用 双 router 模式
其中 internal router 只监听服务器的端口
不通过反向代理进行暴露

---

## USAGE

Use `python memberPyCLI/main.py` for running this shell
Then you can type `help` for menu

| Command | Effect |
| --------------- | --------------- |
| list | get a list of all members |
| create | create a new member |
| update | update a member info |
| delete | delete a member from database |
| query | query member matchs the condition |
| exit | exit the shell |
