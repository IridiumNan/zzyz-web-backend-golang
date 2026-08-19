from typing import Any
from member import Member, MemberConnect, MemberAttr, ZZYZMethod


class MemberShell:
    def __init__(self, base_url: str):
        self.con = MemberConnect(base_url)
        self.menu = """
        list: get a list of all members
        create: create a new member
        update: update a member info
        delete: delete a member from database
        query: query members match the condition
        exit: exit this shell
        """
        self.prompt = "->"

    def ask_with_default(self, attr: str, default_value: Any) -> str:
        val = input(f"set the {attr} [defualt: {default_value}] ->")

        if val == "":
            return default_value

        return val

    def ask_with_required(self, attr: str) -> str:
        val = ""
        while True:
            if val != "":
                return val
            val = input(f"set the {attr} ->")

    def ask_with_optional(self, attr: str) -> str:
        val = input(f"set the {attr} [this is optional just enter for empty] ->")

        return val

    def ask_with_options(self, attr: str, options_list: list[str], default_option: str):
        query = f"set the {attr} ["
        for option in options_list:
            query += option + "/"
        query += "] ->"

        val = input(query)
        if val == "":
            return default_option

        return val

    def do_help(self):
        print(self.menu)

    def do_list(self):
        self.con.new_list_request()

    def do_craete(self):
        print("exec create")
        power = int(self.ask_with_default("power", 0))
        nick = self.ask_with_required("nick")
        email = self.ask_with_optional("email")
        passwd = self.ask_with_required("password")
        is_delete = False

        new_member = Member(
            power=power,
            nick=nick,
            email=email,
            passwd=passwd,
            is_delete=is_delete,
        )

        self.con.new_create_request(new_member)

    def do_update(self):
        # TODO: ask infor interactively
        print("exec update")

    def do_delete(self):
        print("exec delete")
        id = int(self.ask_with_required("id"))

        soft = self.ask_with_options("soft", ["Y", "n"], "y")

        soft_bool = False
        if soft.lower() == "y":
            soft_bool = True

        self.con.new_delete_request(id, soft_bool)

    def do_query(self):
        print("exec queyr")

    def do_exit(self):
        print("Bye")

    def Run(self):
        print("member shell starting, enter help for help manual")
        while True:
            cmd = input(self.prompt)

            match cmd:
                case "help":
                    self.do_help()
                case "list":
                    self.do_list()
                case "create":
                    self.do_craete()
                case "update":
                    self.do_update()
                case "delete":
                    self.do_delete()
                case "query":
                    self.do_query()
                case "exit":
                    self.do_exit()
                    return
                case "":
                    pass
                case _:
                    print("unknown command")
