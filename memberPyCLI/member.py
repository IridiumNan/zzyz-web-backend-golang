from typing import Any
from enum import Enum, auto
import json
import requests
import subprocess


class ZZYZMethod(Enum):
    List = auto()
    Query = auto()
    Create = auto()
    Update = auto()
    Delete = auto()


class MemberAttr(Enum):
    Empty = auto()
    ID = auto()
    Power = auto()
    Nick = auto()
    Email = auto()
    Passwd = auto()
    IsDelete = auto()


def dict_to_json_str(data: dict) -> str:
    return json.dumps(data, indent=4)


def jq_print(data: dict) -> None:

    cmd = ["jq", "-n", "-C", json.dumps(data)]
    result = subprocess.run(cmd, capture_output=True, text=True)

    print(result.stdout)


class Member:
    def __init__(
        self,
        id: int = 0,
        power: int = 0,
        nick: str = "",
        email: str = "",
        passwd: str = "",
        is_delete: bool = False,
    ):
        self.id = id
        self.power = power
        self.nick = nick
        self.email = email
        self.passwd = passwd
        self.is_delete = is_delete

    def set_power(self, power):
        self.power = power

    def set_nick(self, nick):
        self.nick = nick

    def set_email(self, email):
        self.email = email

    def set_passwd(self, passwd):
        self.passwd = passwd

    def set_is_delete(self, is_delete):
        self.is_delete = is_delete

    def get_dict(self):
        return {
            # NOTE: The id just used for mark and update
            # Should not use it as send request
            # "id": self.id,
            "power": self.power,
            "nick": self.nick,
            "email": self.email,
            "passwd": self.passwd,
            "is_delete": self.is_delete,
        }

    def get_json_str(self):
        return json.dumps(self.get_dict())

    def __str__(self) -> str:
        return f"ID: {self.id} | nick: {self.nick} | email: {self.email} | is_delete: {self.is_delete}"


class MemberConnect:
    def __init__(self, base_url: str = "http://127.0.0.1:8081"):
        """
        This base_url not end with slash / !!!
        It just look like http://127.0.0.1:8080
        """

        self.session = requests.Session()

        self.base_url = base_url

        self.endpoints = {
            ZZYZMethod.List: "/member/list",
            ZZYZMethod.Query: "/member/query",
            ZZYZMethod.Create: "/member/create",
            ZZYZMethod.Update: "/member/update",
            ZZYZMethod.Delete: "/member/delete",
        }

    def new_list_request(self):
        url = self.base_url + self.endpoints[ZZYZMethod.List]

        print(f"send a get request to url: {url}")
        resp: requests.Response = self.session.get(url)

        print(f"status_code: {resp.status_code}")

        jq_print(resp.json())

    def new_query_request(
        self,
        attr: MemberAttr = MemberAttr.Empty,
        value: Any = None,
    ) -> None:
        """
        Query with specific attribute
        If select_all == True
        return all Members
        Just use it as you need list all members
        """
        if attr == MemberAttr.Empty:
            return

        url = self.base_url + self.endpoints[ZZYZMethod.Query]

        match attr:
            case MemberAttr.ID:
                url += "/id"
            case MemberAttr.Nick:
                url += "/nick"
                value = f"%{value}%"
            case MemberAttr.Power:
                url += "/power"
            case MemberAttr.Email:
                url += "/email"
                value = f"%{value}%"
            case MemberAttr.IsDelete:
                url += "/is_delete"
        print(f"send a get request to url: {url}, vale: {value}")
        resp: requests.Response = self.session.get(url, params=f"value={value}")

        print("status code: ", resp.status_code)

        jq_print(resp.json())

    def new_create_request(self, m: Member) -> None:
        url = self.base_url + self.endpoints[ZZYZMethod.Create]
        member_dict = m.get_dict()

        payload = {"data": {"member": member_dict}}

        resp: requests.Response = self.session.post(url, json=payload)

        print("status code: ", resp.status_code)

        jq_print(resp.json())

    def new_update_request(self, id: int, attr: str, value: str) -> None:
        url = (
            self.base_url
            + self.endpoints[ZZYZMethod.Update]
            + "/"
            + str(id)
            + "/"
            + attr
        )

        resp: requests.Response = self.session.patch(url, params=f"value={value}")

        print("status code: ", resp.status_code)

        jq_print(resp.json())

    def new_delete_request(self, id: int, soft: bool) -> None:
        url = self.base_url + self.endpoints[ZZYZMethod.Delete] + "/" + str(id)

        print(f"send a request to url: {url}, soft: {soft}")

        resp: requests.Response = self.session.delete(url, params=f"soft={soft}")

        print("status code: ", resp.status_code)

        jq_print(resp.json())
