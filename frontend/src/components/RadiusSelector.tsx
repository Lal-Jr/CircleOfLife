"use client";

import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Compass } from "lucide-react";

interface RadiusSelectorProps {
    value: string;
    onChange: (value: string) => void;
}

export function RadiusSelector({ value, onChange }: RadiusSelectorProps) {
    return (
        <div className="flex items-center gap-2">
            <div className="flex items-center text-sm font-medium text-muted-foreground mr-1">
                <Compass className="h-4 w-4 mr-1.5" />
                Radius:
            </div>
            <Select value={value} onValueChange={onChange}>
                <SelectTrigger className="w-[110px] h-8 text-sm bg-background border-muted shadow-sm hover:bg-accent/50 transition-colors focus:ring-1 focus:ring-primary focus:border-primary">
                    <SelectValue placeholder="Select radius" />
                </SelectTrigger>
                <SelectContent align="end" className="animate-in fade-in-80 slide-in-from-top-1">
                    <SelectItem value="1" className="cursor-pointer focus:bg-primary/10">1 km</SelectItem>
                    <SelectItem value="3" className="cursor-pointer focus:bg-primary/10">3 km</SelectItem>
                    <SelectItem value="5" className="cursor-pointer focus:bg-primary/10">5 km</SelectItem>
                    <SelectItem value="10" className="cursor-pointer focus:bg-primary/10">10 km</SelectItem>
                </SelectContent>
            </Select>
        </div>
    );
}
