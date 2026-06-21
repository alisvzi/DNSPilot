import { useEffect, useState } from "react";

import { GetNetworkAdapters } from "../wailsjs/go/main/App";

function App() {
  const [adapters, setAdapters] = useState<any[]>([]);

  useEffect(() => {
    GetNetworkAdapters().then(setAdapters).catch(console.error);
  }, []);

  return (
    <div>
      <h1>DNSPilot</h1>

      <h2>Network Adapters</h2>

      {adapters.map((adapter) => (
        <div key={adapter.id}>
          <strong>{adapter.name}</strong>

          <div>{adapter.ipv4}</div>
        </div>
      ))}
    </div>
  );
}

export default App;
