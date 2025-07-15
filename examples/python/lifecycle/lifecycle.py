from hanzo_runtime import HanzoRuntime


def main():
    hanzo_runtime = HanzoRuntime()

    print("Creating sandbox")
    sandbox = hanzo_runtime.create()
    print("Sandbox created")

    sandbox.set_labels(
        {
            "public": True,
        }
    )

    print("Stopping sandbox")
    hanzo_runtime.stop(sandbox)
    print("Sandbox stopped")

    print("Starting sandbox")
    hanzo_runtime.start(sandbox)
    print("Sandbox started")

    print("Getting existing sandbox")
    existing_sandbox = hanzo_runtime.get(sandbox.id)
    print("Get existing sandbox")

    response = existing_sandbox.process.exec('echo "Hello World from exec!"', cwd="/home/hanzo", timeout=10)
    if response.exit_code != 0:
        print(f"Error: {response.exit_code} {response.result}")
    else:
        print(response.result)

    sandboxes = hanzo_runtime.list()
    print("Total sandboxes count:", len(sandboxes))

    print(f"Printing sandboxes[0] -> id: {sandboxes[0].id} state: {sandboxes[0].state}")

    print("Removing sandbox")
    hanzo_runtime.delete(sandbox)
    print("Sandbox removed")


if __name__ == "__main__":
    main()
