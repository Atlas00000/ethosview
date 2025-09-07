"use client";
import React, { useEffect, useRef } from "react";

export function ScrollReveal({ children }: { children: React.ReactNode }) {
  const ref = useRef<HTMLDivElement | null>(null);
  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const obs = new IntersectionObserver((entries) => {
      entries.forEach((e) => {
        if (e.isIntersecting) {
          e.target.classList.add("reveal-visible");
          obs.unobserve(e.target);
        }
      });
    }, { rootMargin: "0px 0px -10% 0px", threshold: 0.15 });
    el.querySelectorAll<HTMLElement>(".reveal-on-scroll").forEach((n) => obs.observe(n));
    return () => obs.disconnect();
  }, []);
  return <div ref={ref}>{children}</div>;
}


