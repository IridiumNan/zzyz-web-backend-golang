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
        self.prompt = "zzyz>"

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
        query = query[:-1] + f"] default {default_option} ->"

        val = ""
        val = input(query)
        if val == "":
            return default_option
        while val not in options_list:
            val = input("invalid choice, select again ->")

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
        attr_list = ["id", "nick", "power", "email", "is_delete"]
        attr = self.ask_with_options("attribute", attr_list, "id")

        match attr:
            case "id":
                attr_val = self.ask_with_required(attr)
                self.con.new_query_request(MemberAttr.ID, attr_val)
            case "nick":
                attr_val = self.ask_with_required(attr)
                self.con.new_query_request(MemberAttr.Nick, attr_val)
            case "power":
                attr_val = self.ask_with_options(attr, ["0", "1"], "0")
                self.con.new_query_request(MemberAttr.Power, attr_val)
            case "email":
                attr_val = self.ask_with_required(attr)
                self.con.new_query_request(MemberAttr.Email, attr_val)
            case "is_delete":
                attr_val = self.ask_with_options(attr, ["0", "1"], "0")
                self.con.new_query_request(MemberAttr.IsDelete, attr_val)

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
