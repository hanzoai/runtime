import { Global, Inject, Module, OnApplicationShutdown } from '@nestjs/common'
import KV from '@hanzo/kv'

import { TypedConfigService } from '../config/typed-config.service'

export const KV_CLIENT = 'KV_CLIENT'

/** Injects the shared KV connection. */
export const InjectKV = () => Inject(KV_CLIENT)

/**
 * One KV connection for the process, provided directly rather than through a
 * wrapper package. The wrapper this replaces existed to construct a client and
 * hand it to Nest's container, which is what the factory below does — and it
 * carried its own copy of a different client, so the process ended up with two.
 */
@Global()
@Module({
  providers: [
    {
      provide: KV_CLIENT,
      inject: [TypedConfigService],
      useFactory: (config: TypedConfigService) =>
        new KV({
          host: config.getOrThrow('kv.host'),
          port: config.getOrThrow('kv.port'),
          tls: config.get('kv.tls'),
          lazyConnect: config.get('skipConnections'),
        }),
    },
  ],
  exports: [KV_CLIENT],
})
export class KVModule implements OnApplicationShutdown {
  constructor(@InjectKV() private readonly kv: KV) {}

  /** Nest owns the connection's lifetime, so it closes with the application. */
  async onApplicationShutdown(): Promise<void> {
    await this.kv.quit()
  }
}
