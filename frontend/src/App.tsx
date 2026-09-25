import { useState } from 'react'
import Navbar from './components/Navbar'
import Drawer from './components/Drawer'
import Calculator from './components/Calculator'
import Footer from './components/Footer'
import './App.css'

function App() {
  const [drawerOpen, setDrawerOpen] = useState(false)

  return (
    <>
      <Navbar onMenuClick={() => setDrawerOpen(true)} />
      <Drawer open={drawerOpen} onClose={() => setDrawerOpen(false)} />

      <main className="page-main">
        <Calculator />
      </main>

      <Footer />
    </>
  )
}

export default App
