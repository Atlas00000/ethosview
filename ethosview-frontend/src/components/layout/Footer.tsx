import React from "react";

export function Footer() {
  return (
    <footer className="mt-16 border-t" style={{ borderColor: "rgba(0,0,0,0.06)" }}>
      <div className="max-w-6xl mx-auto px-4 py-8 text-sm flex flex-col sm:flex-row items-center justify-between gap-3">
        <div className="text-gradient font-semibold">EthosView</div>
        <div className="text-gray-600">ESG and Financial Analytics</div>
        <div className="text-gray-500">© {new Date().getFullYear()} EthosView</div>
      </div>
    </footer>
  );
}


