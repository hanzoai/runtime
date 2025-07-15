import { HanzoRuntime } from '@hanzo/runtime-sdk'
import path from 'path'

async function main() {
  const hanzoRuntime = new HanzoRuntime()

  //  Create a new volume or get an existing one
  const volume = await hanzoRuntime.volume.get('my-volume', true)

  // Mount the volume to the sandbox
  const mountDir1 = '/home/hanzo/volume'

  const sandbox1 = await hanzoRuntime.create({
    language: 'typescript',
    volumes: [{ volumeId: volume.id, mountPath: mountDir1 }],
  })

  // Create a new directory in the mount directory
  const newDir = path.join(mountDir1, 'new-dir')
  await sandbox1.fs.createFolder(newDir, '755')

  // Create a new file in the mount directory
  const newFile = path.join(mountDir1, 'new-file.txt')
  await sandbox1.fs.uploadFile(Buffer.from('Hello, World!'), newFile)

  // Create a new sandbox with the same volume
  // and mount it to the different path
  const mountDir2 = '/home/hanzo/my-files'

  const sandbox2 = await hanzoRuntime.create({
    language: 'typescript',
    volumes: [{ volumeId: volume.id, mountPath: mountDir2 }],
  })

  // List files in the mount directory
  const files = await sandbox2.fs.listFiles(mountDir2)
  console.log('Files:', files)

  // Get the file from the first sandbox
  const file = await sandbox1.fs.downloadFile(newFile)
  console.log('File:', file.toString())

  // Cleanup
  await hanzoRuntime.delete(sandbox1)
  await hanzoRuntime.delete(sandbox2)
  // await hanzoRuntime.volume.delete(volume)
}

main()
