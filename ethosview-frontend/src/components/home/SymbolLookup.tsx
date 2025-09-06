"use client";
import React, { useState } from "react";
import { api } from "../../services/api";
import type { CompanyResponse } from "../../types/api";

export function SymbolLookup() {
  const [symbol, setSymbol] = useState("");
  const [result, setResult] = useState<CompanyResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSearch(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setResult(null);
    if (!symbol.trim()) return;
    setLoading(true);
    try {
      const company = await api.companyBySymbol(symbol.trim());
      setResult(company as CompanyResponse);
    } catch (err: any) {
      setError("Not found");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="max-w-6xl mx-auto px-4 py-6 animate-fade-in-up">
      <form onSubmit={onSearch} className="flex gap-2">
        <input
          className="border rounded px-3 py-2 w-64 focus:outline-none focus:ring-2"
          placeholder="Lookup by symbol (e.g., AAPL)"
          value={symbol}
          onChange={(e) => setSymbol(e.target.value.toUpperCase())}
        />
        <button className="btn-primary rounded px-3 py-2" disabled={loading}>
          {loading ? "Searching..." : "Search"}
        </button>
      </form>
      {error && <div className="text-sm text-red-600 mt-2">{error}</div>}
      {result && (
        <div className="mt-3 text-sm">
          <div className="font-medium">{result.name} ({result.symbol})</div>
          <div className="text-gray-600">Sector: {result.sector}</div>
        </div>
      )}
    </section>
  );
}


