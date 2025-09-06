export default function RootLoading() {
  return (
    <div className="min-h-[50vh] flex items-center justify-center">
      <div className="glass-card p-6 shadow-elevated">
        <div className="w-64 h-2 rounded-full skeleton" />
        <div className="mt-3 text-gradient font-semibold">Loading EthosView…</div>
      </div>
    </div>
  );
}


