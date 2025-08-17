/*
 * Copyright 2025 Hanzo Industries Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

export interface OrganizationWallet {
  balanceCents: number
  ongoingBalanceCents: number
  name: string
  creditCardConnected: boolean

  automaticTopUp?: AutomaticTopUp
}

export type AutomaticTopUp = {
  thresholdAmount: number
  targetAmount: number
}
