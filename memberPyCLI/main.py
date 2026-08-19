from shell import MemberShell

base_url = "http://127.0.0.1:8081"


if __name__ == "__main__":
    sh = MemberShell(base_url)
    sh.Run()
