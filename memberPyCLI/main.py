from member import Member, MemberConnect, MemberAttr, ZZYZMethod

base_url = "http://127.0.0.1:8081"

con = MemberConnect(base_url)

json_resp = con.new_delete_request(1)

print(json_resp)
