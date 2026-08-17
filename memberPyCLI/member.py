from typing import Any
from enum import Enum, auto
import json
import requests


class ZZYZMethod(Enum):
    Query = auto()
    Create = auto()
    Update = auto()
    Delete = auto()


class MemberAttr(Enum):
    Empty = auto()
    ID = auto()
    Nick = auto()
    Email = auto()
    Passwd = auto()
    IsDelete = auto()


class Member:
    def __init__(
        self,
        id: int = -1,
        nick: str = "",
        email: str = "",
        passwd: str = "",
        is_delete: bool = False,
    ):
        self.id = id
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
            # Don't send it when request
            # "id": self.id,
            "nick": self.nick,
            "email": self.email,
            "password": self.passwd,
            "is_delete": self.is_delete,
        }

    def get_json_str(self):
        return json.dumps(self.get_dict())

    def __str__(self) -> str:
        return f"ID: {self.id} | nick: {self.nick} | email: {self.email} | is_delete: {self.is_delete}"


class MemberConnect:
    def __init__(self, base_url: str):
        """
        This base_url not end with slash / !!!
        It just look like http://127.0.0.1:8080
        """

        self.session = requests.Session()

        self.base_url = base_url

        self.endpoints = {
            ZZYZMethod.Query: "/member/query",
            ZZYZMethod.Create: "/member/create",
            ZZYZMethod.Update: "/member/update",
            ZZYZMethod.Delete: "/member/delete",
        }

    def new_query_request(
        self,
        member_attr: MemberAttr = MemberAttr.Empty,
        value: Any = None,
        select_all: bool = False,
    ) -> Any:
        """
        Query with specific attributte
        If select_all == True
        return all Members
        Just use it as you need list all members
        """

        url = self.base_url + self.endpoints[ZZYZMethod.Query]
        pass

    def new_create_request(self, m: Member) -> Any:
        url = self.base_url + self.endpoints[ZZYZMethod.Create]
        member_dict = m.get_dict()

        payload = {"data": {"member": member_dict}}

        resp: requests.Response = self.session.post(url, json=payload)

        print("status code: ", resp.status_code)

        return resp.json()

    def new_update_request(self, id: int, updated_member: Member) -> Any:
        url = self.base_url + self.endpoints[ZZYZMethod.Update]
        member_dict = updated_member.get_dict()

        payload = {"data": {"id": id, "member": member_dict}}

        resp: requests.Response = self.session.patch(url, json=payload)

        print("status code: ", resp.status_code)

        return resp.json()

    def new_delete_request(self, id: int) -> Any:
        url = self.base_url + self.endpoints[ZZYZMethod.Delete]

        print(f"send a request to url: {url}, id: {id}")

        resp: requests.Response = self.session.delete(url, params=f"id={id}")

        print("status code: ", resp.status_code)

        return resp.json()
