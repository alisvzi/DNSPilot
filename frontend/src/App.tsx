import { useEffect, useState } from "react";
import { GetSystemDNS } from "../wailsjs/go/main/App";

function App() {
  const [data, setData] = useState<any[]>([]);

  useEffect(() => {
    GetSystemDNS().then(setData);
  }, []);

  return (
    <div style={{ padding: 20 }}>
      <h1>DNSPilot</h1>

      {data.map((item, i) => (
        <div key={i} style={{ marginBottom: 20 }}>
          <h3>{item.adapter_name}</h3>

          {item.dns_servers.map((dns: string, j: number) => (
            <div key={j}>• {dns}</div>
          ))}
        </div>
      ))}
    </div>
  );
}

export default App;
