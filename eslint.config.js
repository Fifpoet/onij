import js from '@eslint/js'
import ts from 'typescript-eslint'
import vue from 'eslint-plugin-vue'
import globals from 'globals'

/**
 * 只拦「真实的错误」，不参与格式之争（缩进/引号交给编辑器）。
 * 格式类规则一开就会产出几千条噪音，反而没人看。
 */
export default ts.config(
  {
    ignores: [
      'dist/**',
      'dist-verify/**',
      'release/**',
      'node_modules/**',
      'stats.html',
      'src/vendor/**',
      'uvr/**',
      'chrome/**',
      'server/**',
    ],
  },
  js.configs.recommended,
  ...ts.configs.recommended,
  ...vue.configs['flat/recommended'],
  {
    files: ['**/*.vue'],
    languageOptions: {
      parserOptions: {
        // .vue 里用 TS 语法必须显式指定解析器
        parser: ts.parser,
        ecmaVersion: 'latest',
        sourceType: 'module',
      },
    },
  },
  {
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
    rules: {
      // 未使用变量：允许下划线前缀（回调参数常需要占位）
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_', caughtErrors: 'none' },
      ],
      // 现存代码里 any 较多，先降为提醒，逐步收
      '@typescript-eslint/no-explicit-any': 'warn',
      '@typescript-eslint/no-empty-object-type': 'off',
      // Vue：组件名允许单词（现有页面都是单文件页面级组件）
      'vue/multi-word-component-names': 'off',
      // 模板格式类规则全部关闭，避免噪音
      'vue/max-attributes-per-line': 'off',
      'vue/singleline-html-element-content-newline': 'off',
      'vue/html-self-closing': 'off',
      'vue/html-indent': 'off',
      'vue/html-closing-bracket-newline': 'off',
      'vue/attributes-order': 'off',
      'vue/first-attribute-linebreak': 'off',
    },
  },
)
