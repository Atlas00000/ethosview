import { createClient } from '@supabase/supabase-js'

const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL!
const supabaseAnonKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!

export const supabase = createClient(supabaseUrl, supabaseAnonKey)

// Database types for TypeScript
export interface Database {
  public: {
    Tables: {
      users: {
        Row: {
          id: string
          email: string
          password_hash: string
          created_at: string
          updated_at: string
        }
        Insert: {
          id?: string
          email: string
          password_hash: string
          created_at?: string
          updated_at?: string
        }
        Update: {
          id?: string
          email?: string
          password_hash?: string
          created_at?: string
          updated_at?: string
        }
      }
      companies: {
        Row: {
          id: string
          name: string
          symbol: string
          sector: string
          created_at: string
        }
        Insert: {
          id?: string
          name: string
          symbol: string
          sector: string
          created_at?: string
        }
        Update: {
          id?: string
          name?: string
          symbol?: string
          sector?: string
          created_at?: string
        }
      }
      esg_scores: {
        Row: {
          id: string
          company_id: string
          environmental_score: number
          social_score: number
          governance_score: number
          overall_score: number
          date: string
          created_at: string
        }
        Insert: {
          id?: string
          company_id: string
          environmental_score: number
          social_score: number
          governance_score: number
          overall_score: number
          date: string
          created_at?: string
        }
        Update: {
          id?: string
          company_id?: string
          environmental_score?: number
          social_score?: number
          governance_score?: number
          overall_score?: number
          date?: string
          created_at?: string
        }
      }
    }
  }
}
