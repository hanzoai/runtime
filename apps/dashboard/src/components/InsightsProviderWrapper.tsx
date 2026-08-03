/*
 * Copyright 2025 Daytona Platforms Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

import { FC, ReactNode } from 'react'
import { InsightsProvider } from '@hanzo/insights/react'

const insightsKey = import.meta.env.VITE_INSIGHTS_KEY
const insightsHost = import.meta.env.VITE_INSIGHTS_HOST

interface InsightsProviderWrapperProps {
  children: ReactNode
}

export const InsightsProviderWrapper: FC<InsightsProviderWrapperProps> = ({ children }) => {
  if (!import.meta.env.PROD) {
    return children
  }

  if (!insightsKey || !insightsHost) {
    console.error('Invalid Insights configuration')
    return children
  }

  return (
    <InsightsProvider
      apiKey={insightsKey}
      options={{
        api_host: insightsHost,
        person_profiles: 'always',
        autocapture: false, // ignore default frontend events
        capture_pageview: false, // initial pageview (handled in App.tsx)
        capture_pageleave: true, // end of session
      }}
    >
      {children}
    </InsightsProvider>
  )
}
