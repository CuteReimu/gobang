package com.github.reimu.gobang;

import java.util.HashMap;
import java.util.Map;

public class BoardCache {
	private final Map<Long, PointAndValueAndDeep> cache;
	
	public BoardCache() {
		cache = new HashMap<>();
	}
	
	/**
	 * 如果存的deep更深，则不put并返回存的Point和value，否则put并返回null
	 */
	public void put(long key, PointAndValue pointAndValue, int deep) {
		put(key, pointAndValue.point, pointAndValue.value, deep);
	}
	
	/**
	 * 如果存的deep更深，则不put并返回存的Point和value，否则put并返回null
	 */
	public void put(long key, Point point, int value, int deep) {
		PointAndValueAndDeep cacheValue = cache.get(key);
		if (null != cacheValue && deep < cacheValue.deep) {
			return;
		}
		cache.put(key, new PointAndValueAndDeep(point, value, deep));
	}
	
	/**
	 * 如果存的deep满足深度要求了，则返回。如果深度不达标或者没存，则返回null
	 */
	public PointAndValue get(long key, int deep) {
		PointAndValueAndDeep cacheValue = cache.get(key);
		if (null != cacheValue && deep <= cacheValue.deep)
			return new PointAndValue(cacheValue.point, cacheValue.value);
		return null;
	}

	private static class PointAndValueAndDeep {
		public final Point point;
		public final int value;
		public final int deep;
		public PointAndValueAndDeep(Point point, int value, int deep) {
			super();
			this.point = point;
			this.value = value;
			this.deep = deep;
		}
	}
}
