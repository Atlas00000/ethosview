export default function RootLoading() {
  return (
    <div className="min-h-[60vh] flex items-center justify-center relative overflow-hidden">
      <div className="blob blob-drift-a" style={{ top: -60, left: -80, width: 220, height: 220, background: "rgba(30,106,225,0.12)", borderRadius: 9999 }} />
      <div className="blob blob-drift-b" style={{ bottom: -60, right: -80, width: 220, height: 220, background: "rgba(42,179,166,0.12)", borderRadius: 9999 }} />
      <div className="max-w-4xl mx-auto px-4 w-full">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="glass-card p-4 animate-fade-in-up">
            <div className="w-28 h-3 rounded skeleton" />
            <div className="mt-3 h-24 rounded skeleton" />
          </div>
          <div className="glass-card p-4 animate-fade-in-up" style={{ animationDelay: "80ms" }}>
            <div className="w-24 h-3 rounded skeleton" />
            <div className="mt-3 space-y-2">
              <div className="h-3 rounded skeleton" />
              <div className="h-3 rounded skeleton" />
              <div className="h-3 rounded skeleton" />
            </div>
          </div>
          <div className="glass-card p-4 animate-fade-in-up" style={{ animationDelay: "140ms" }}>
            <div className="w-20 h-3 rounded skeleton" />
            <div className="mt-3 grid grid-cols-3 gap-2">
              <div className="h-12 rounded skeleton" />
              <div className="h-12 rounded skeleton" />
              <div className="h-12 rounded skeleton" />
            </div>
          </div>
        </div>
        <div className="mt-6 glass-card p-4 animate-fade-in-up" style={{ animationDelay: "200ms" }}>
          <div className="w-36 h-3 rounded skeleton" />
          <div className="mt-3 h-2 rounded bg-white/40 overflow-hidden">
            <div className="h-2 w-1/3 bg-gradient-to-r from-[#1E6AE1] to-[#2AB3A6] animate-[loading_1.2s_ease-in-out_infinite]" />
          </div>
          <div className="mt-2 text-xs" style={{ color: "#374151" }}>Preparing live ESG and market data…</div>
        </div>
        <div className="mt-4 flex items-center justify-center">
          <div className="dot-pulse" />
          <div className="ml-2 text-gradient font-semibold">Loading EthosView…</div>
        </div>
      </div>
    </div>
  );
}


