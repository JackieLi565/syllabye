export interface Session {
  userId: string
}

export interface User {
  id: string
  currentYear?: number | null
  email: string
  fullname: string
  gender?: 'male' | 'female' | 'other'
  nickname?: string
  picture: string
  programId?: string
}