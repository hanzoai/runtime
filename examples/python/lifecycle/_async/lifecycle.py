import asyncio

from hanzo_runtime import AsyncHanzoRuntime


async def main():
    async with AsyncHanzoRuntime() as hanzo_runtime:
        print("Creating sandbox")
        sandbox = await hanzo_runtime.create()
        print("Sandbox created")

        await sandbox.set_labels(
            {
                "public": True,
            }
        )

        print("Stopping sandbox")
        await hanzo_runtime.stop(sandbox)
        print("Sandbox stopped")

        print("Starting sandbox")
        await hanzo_runtime.start(sandbox)
        print("Sandbox started")

        print("Getting existing sandbox")
        existing_sandbox = await hanzo_runtime.get(sandbox.id)
        print("Get existing sandbox")

        response = await existing_sandbox.process.exec('echo "Hello World from exec!"', cwd="/home/hanzo", timeout=10)
        if response.exit_code != 0:
            print(f"Error: {response.exit_code} {response.result}")
        else:
            print(response.result)

        sandboxes = await hanzo_runtime.list()
        print("Total sandboxes count:", len(sandboxes))

        print(f"Printing sandboxes[0] -> id: {sandboxes[0].id} state: {sandboxes[0].state}")

        print("Removing sandbox")
        await hanzo_runtime.delete(sandbox)
        print("Sandbox removed")


if __name__ == "__main__":
    asyncio.run(main())
