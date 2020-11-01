import os
import re
import shutil

def checkIP(ip):
    _ip_regex = re.compile("^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$")
    m = _ip_regex.match(ip)
    if m != None:
        return True
    else:
        return False

def enter_details(node_count):
    dictionary = {}
    for x in range(node_count):
        if node_count > 1:
            print(format(x+1) + "/" + format(node_count))
        while True:
            name = input("Name > ")
            if not name in dictionary and name != "":
                break
            elif name == "":
                print("Name can't be blank. Try again.")
            else:
                print("You already got that node. Try again.")
        while True:
            ip = input("IPv4 > ")
            if checkIP(ip) and not ip in dictionary.values():
                break
            elif ip in dictionary.values():
                print("You already used that IP. Try again.")
            else:
                print("That looks like a malformed IP. Try again.")
        dictionary[name] = ip
    return dictionary

def getlines(fobj,line1,line2):
    for line in iter(fobj.readline,''):  #This is necessary to get `fobj.tell` to work
        yield line
        if line == line1:
            pos = fobj.tell()
            next_line = next(fobj)
            fobj.seek(pos)
            if next_line == line2:
                return

def main():

    bootstrap = {}
    masters = {}
    workers = {}

    print("Let's generate an HAProxy config.")

    print("First up, what are the deets on your bootstrap?")
    bootstrap = enter_details(1)

    print("Okay, how many masters do you want?")

    master_count = int(input(" > "))

    masters_as_workers = False
    check = ""
    while True:
        check = input("Are they going to be used as workers? (y/n) ")
        print(check)
        if check == "y" or check == "yes":
            masters_as_workers = True
            break
        elif check == "n" or check == "no":
            break
        
    print("Enter your master details")
    masters = enter_details(master_count)

    print("Okay, how many workers do you want?")

    worker_count = int(input(" > "))

    if worker_count > 0:
        print("Enter your worker details")
        workers = enter_details(worker_count)
    
    print("Generating config...")

    try:
        os.mkdir("../okd4_output")
        
    except:
        pass
    
    try:
        os.remove("../okd4_output/haproxy.conf")
    except:
        pass

    with open('../template/haproxy.conf') as fin, open('../okd4_output/haproxy.conf','w+') as fout:
        # Write bootstrap and masters
        fout.writelines(getlines(fin,'    mode tcp\n','    # bootstrap and masters go here\n'))
        fout.write("    server      {} {}:6443 check\n".format(list(bootstrap.keys())[0], list(bootstrap.values())[0]))
        for i in range(master_count):
            fout.write("    server      {} {}:6443 check\n".format(list(masters.keys())[0], list(masters.values())[0]))

        fout.writelines(getlines(fin,'    mode tcp\n','    # bootstrap and masters go here\n'))
        fout.write("    server      {} {}:22623 check\n".format(list(bootstrap.keys())[0], list(bootstrap.values())[0]))
        for i in range(master_count):
            fout.write("    server      {} {}:22623 check\n".format(list(masters.keys())[0], list(masters.values())[0]))

        # Write the workers
        fout.writelines(getlines(fin,'    mode tcp\n','    # computes go here\n'))
        for i in range(worker_count):
            fout.write("    server      {} {}:80 check\n".format(list(workers.keys())[i], list(workers.values())[i]))

        fout.writelines(getlines(fin,'    mode tcp\n','    # computes go here\n'))
        for i in range(worker_count):
            fout.write("    server      {} {}:443 check\n".format(list(workers.keys())[i], list(workers.values())[i]))
        fout.writelines(fin)
    print("Done.")

if __name__ == "__main__":
    main()