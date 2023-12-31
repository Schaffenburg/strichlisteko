/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./public/**/*.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'media', // or 'class'
  theme: {
    extend:{
      colors:{
        'bunker': {  DEFAULT: '#060A0D',  50: '#1C2E3C',  100: '#192A37',  200: '#14222C',  300: '#101A22',  400: '#0B1217',  500: '#060A0D',  600: '#04070A',  700: '#030506',  800: '#010203',  900: '#000000',  950: '#000000'},
        'lgray': {  DEFAULT: '#434F59',  50: '#9EABB6',  100: '#92A1AD',  200: '#7B8D9C',  300: '#667888',  400: '#556470',  500: '#434F59',  600: '#364048',  700: '#293036',  800: '#1C2125',  900: '#0E1113',  950: '#08090A'},
        'concrete': {  DEFAULT: '#F2F2F2',  50: '#FDFDFD',  100: '#FCFCFC',  200: '#FAFAFA',  300: '#F7F7F7',  400: '#F5F5F5',  500: '#F2F2F2',  600: '#D6D6D6',  700: '#BABABA',  800: '#9E9E9E',  900: '#828282',  950: '#747474'},
        'primary': {  DEFAULT: '#1B418C',  50: '#799EE5',  100: '#6891E2',  200: '#4678DB',  300: '#2861D0',  400: '#2251AE',  500: '#1B418C',  600: '#15336E',  700: '#0F2550',  800: '#0A1732',  900: '#040914',  950: '#010205'},
      }
    }
  },
  plugins: [],
}

