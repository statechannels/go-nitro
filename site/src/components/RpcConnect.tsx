import { useState } from "react";

export default function RpcConnect() {
  const [text, setText] = useState("localhost:4005");

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setText(e.target.value);
  }

  return (
    <div>
      RPC Connect
      <input value={text} onChange={handleChange} />
      <button
        onClick={() => {
          console.log("Button click");
        }}
      >
        Connect
      </button>
    </div>
  );
}
