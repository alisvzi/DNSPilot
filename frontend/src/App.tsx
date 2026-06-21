import { useEffect, useState } from "react";
import { GetSystemDNS } from "../wailsjs/go/main/App";

function App() {
  const [dns, setDns] = useState<string[]>([]);

  useEffect(() => {
    GetSystemDNS().then(setDns).catch(console.error);
  }, []);

  return (
    <div style={{ padding: 20 }}>
      <h1>DNSPilot</h1>

      <h2>System DNS</h2>

      {dns.map((d, i) => (
        <div key={i}>{d}</div>
      ))}
    </div>
  );
}

export default App;
