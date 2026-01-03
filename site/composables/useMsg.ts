export const useMsgSuccess = (message: string) => {
  console.log('success:', message)
}

export const useMsgError = (message: string, error?: unknown) => {
  console.log('error:', message, error)
}
