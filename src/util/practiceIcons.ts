import type { Component } from 'vue'
import {
  BarbellOutline,
  BicycleOutline,
  BookOutline,
  EllipsisHorizontalOutline,
  GameControllerOutline,
  HelpCircleOutline,
  LanguageOutline,
  MicOutline,
  MusicalNoteOutline,
  MusicalNotesOutline,
} from '@vicons/ionicons5'
import { PracticeType } from '@/api/types/enums'

/** 各练习类型对应图标 */
export const PRACTICE_TYPE_ICONS: Record<PracticeType, Component> = {
  [PracticeType.PCT_Unknown]: HelpCircleOutline,
  [PracticeType.PCT_AnaerobicExercise]: BarbellOutline,
  [PracticeType.PCT_AerobicExercise]: BicycleOutline,
  [PracticeType.PCT_Reading]: BookOutline,
  [PracticeType.PCT_Language]: LanguageOutline,
  [PracticeType.PCT_Game]: GameControllerOutline,
  [PracticeType.PCT_Piano]: MusicalNotesOutline,
  [PracticeType.PCT_Flute]: MusicalNoteOutline,
  [PracticeType.PCT_Singing]: MicOutline,
  [PracticeType.PCT_Other]: EllipsisHorizontalOutline,
}

export function practiceIcon(type: PracticeType): Component {
  return PRACTICE_TYPE_ICONS[type] ?? PRACTICE_TYPE_ICONS[PracticeType.PCT_Other]
}
