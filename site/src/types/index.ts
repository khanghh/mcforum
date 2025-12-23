
export type CursorList<T> = {
  url: string
  items: T[]
  hasMore: boolean
  cursor: number
}

export * from './user'
export * from './topic'