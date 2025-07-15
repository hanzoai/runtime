import os

from hanzo_runtime import , VolumeMount


def main():
    hanzo_runtime = HanzoRuntime()

    # Create a new volume or get an existing one
    volume = hanzo_runtime.volume.get("my-volume", create=True)

    # Mount the volume to the sandbox
    mount_dir_1 = "/home/hanzo/volume"

    params = CreateSandboxFromSnapshotParams(
        language="python",
        volumes=[VolumeMount(volumeId=volume.id, mountPath=mount_dir_1)],
    )
    sandbox = hanzo_runtime.create(params)

    # Create a new directory in the mount directory
    new_dir = os.path.join(mount_dir_1, "new-dir")
    sandbox.fs.create_folder(new_dir, "755")

    # Create a new file in the mount directory
    new_file = os.path.join(mount_dir_1, "new-file.txt")
    sandbox.fs.upload_file(b"Hello, World!", new_file)

    # Create a new sandbox with the same volume
    # and mount it to the different path
    mount_dir_2 = "/home/hanzo/my-files"

    params = CreateSandboxFromSnapshotParams(
        language="python",
        volumes=[VolumeMount(volumeId=volume.id, mountPath=mount_dir_2)],
    )
    sandbox2 = hanzo_runtime.create(params)

    # List files in the mount directory
    files = sandbox2.fs.list_files(mount_dir_2)
    print("Files:", files)

    # Get the file from the mount directory
    file = sandbox2.fs.download_file(os.path.join(mount_dir_2, "new-file.txt"))
    print("File:", file)

    # Cleanup
    hanzo_runtime.delete(sandbox)
    hanzo_runtime.delete(sandbox2)
    # hanzo_runtime.volume.delete(volume)


if __name__ == "__main__":
    main()
