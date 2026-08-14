/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { MigrationInterface, QueryRunner } from 'typeorm'

// Every organization and every sandbox that predates this column reads as
// gvisor, so existing rows come out isolated rather than needing a backfill to
// become safe.
export class Migration1755100000000 implements MigrationInterface {
  name = 'Migration1755100000000'

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`CREATE TYPE "public"."isolation_enum" AS ENUM('gvisor', 'firecracker', 'runc')`)
    await queryRunner.query(
      `ALTER TABLE "organization" ADD "sandbox_isolation" "public"."isolation_enum" NOT NULL DEFAULT 'gvisor'`,
    )
    await queryRunner.query(
      `ALTER TABLE "sandbox" ADD "isolation" "public"."isolation_enum" NOT NULL DEFAULT 'gvisor'`,
    )
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "sandbox" DROP COLUMN "isolation"`)
    await queryRunner.query(`ALTER TABLE "organization" DROP COLUMN "sandbox_isolation"`)
    await queryRunner.query(`DROP TYPE "public"."isolation_enum"`)
  }
}
